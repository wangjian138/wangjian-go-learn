package learn

import (
	"errors"
	"fmt"
	"os"
	"path"
	"shorturl/wangjian-zero/core/logx"
	"strconv"
	"strings"
	"time"

	"shorturl/wangjian-zero/core/iox"
	"shorturl/wangjian-zero/core/lang"
)

const cgroupDir = "/sys/fs/cgroup"

type cgroup struct {
	cgroups map[string]string
}

func (c *cgroup) acctUsageAllCpus() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage"))
	if err != nil {
		return 0, err
	}

	return parseUint(string(data))
}

func (c *cgroup) acctUsagePerCpu() ([]uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuacct"], "cpuacct.usage_percpu"))
	if err != nil {
		return nil, err
	}

	var usage []uint64
	for _, v := range strings.Fields(string(data)) {
		u, err := parseUint(v)
		if err != nil {
			return nil, err
		}

		usage = append(usage, u)
	}

	return usage, nil
}

func (c *cgroup) cpuQuotaUs() (int64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_quota_us"))
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(data), 10, 64)
}

func (c *cgroup) cpuPeriodUs() (uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpu"], "cpu.cfs_period_us"))
	if err != nil {
		return 0, err
	}

	return parseUint(string(data))
}

func (c *cgroup) cpus() ([]uint64, error) {
	data, err := iox.ReadText(path.Join(c.cgroups["cpuset"], "cpuset.cpus"))
	if err != nil {
		return nil, err
	}

	return parseUints(string(data))
}

func currentCgroup() (*cgroup, error) {
	cgroupFile := fmt.Sprintf("/proc/%d/cgroup", os.Getpid())
	lines, err := iox.ReadTextLines(cgroupFile, iox.WithoutBlank())
	if err != nil {
		return nil, err
	}

	cgroups := make(map[string]string)
	for _, line := range lines {
		cols := strings.Split(line, ":")
		if len(cols) != 3 {
			return nil, fmt.Errorf("invalid cgroup line: %s", line)
		}

		subsys := cols[1]
		// only read cpu staff
		if !strings.HasPrefix(subsys, "cpu") {
			continue
		}

		fmt.Printf("currentCgroup line:%v cols:%v\n", line, cols)
		// https://man7.org/linux/man-pages/man7/cgroups.7.html
		// comma-separated list of controllers for cgroup version 1
		fields := strings.Split(subsys, ",")
		for _, val := range fields {
			cgroups[val] = path.Join(cgroupDir, val)
		}
	}

	return &cgroup{
		cgroups: cgroups,
	}, nil
}

func parseUint(s string) (uint64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if err.(*strconv.NumError).Err == strconv.ErrRange {
			return 0, nil
		}

		return 0, fmt.Errorf("cgroup: bad int format: %s", s)
	}

	if v < 0 {
		return 0, nil
	}

	return uint64(v), nil
}

func parseUints(val string) ([]uint64, error) {
	if val == "" {
		return nil, nil
	}

	ints := make(map[uint64]lang.PlaceholderType)
	cols := strings.Split(val, ",")
	for _, r := range cols {
		if strings.Contains(r, "-") {
			fields := strings.SplitN(r, "-", 2)
			min, err := parseUint(fields[0])
			if err != nil {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			max, err := parseUint(fields[1])
			if err != nil {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			if max < min {
				return nil, fmt.Errorf("cgroup: bad int list format: %s", val)
			}

			for i := min; i <= max; i++ {
				ints[i] = lang.Placeholder
			}
		} else {
			v, err := parseUint(r)
			if err != nil {
				return nil, err
			}

			ints[v] = lang.Placeholder
		}
	}

	var sets []uint64
	for k := range ints {
		sets = append(sets, k)
	}

	return sets, nil
}

const (
	cpuTicks  = 100
	cpuFields = 8
)

var (
	preSystem uint64
	preTotal  uint64
	quota     float64
	cores     uint64
)

// if /proc not present, ignore the cpu calcuation, like wsl linux
func Init() {
	cpus, err := perCpuUsage()

	if err != nil {
		logx.Error(err)
		return
	}

	cores = uint64(len(cpus))
	sets, err := cpuSets()
	if err != nil {
		logx.Error(err)
		return
	}

	quota = float64(len(sets))
	cq, err := cpuQuota()
	if err == nil {
		if cq != -1 {
			period, err := cpuPeriod()
			if err != nil {
				logx.Error(err)
				return
			}

			limit := float64(cq) / float64(period)
			if limit < quota {
				quota = limit
			}
		}
	}

	preSystem, err = systemCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}

	preTotal, err = totalCpuUsage()
	if err != nil {
		logx.Error(err)
		return
	}
}

// RefreshCpu refreshes cpu usage and returns.
func RefreshCpu() uint64 {
	total, err := totalCpuUsage()
	if err != nil {
		return 0
	}
	system, err := systemCpuUsage()
	if err != nil {
		return 0
	}

	var usage uint64
	cpuDelta := total - preTotal
	systemDelta := system - preSystem
	if cpuDelta > 0 && systemDelta > 0 {
		usage = uint64(float64(cpuDelta*cores*1e3) / (float64(systemDelta) * quota))
	}
	preSystem = system
	preTotal = total

	return usage
}

func cpuQuota() (int64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuQuotaUs()
}

func cpuPeriod() (uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuPeriodUs()
}

func cpuSets() ([]uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return nil, err
	}

	return cg.cpus()
}

func perCpuUsage() ([]uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return nil, err
	}

	return cg.acctUsagePerCpu()
}

func systemCpuUsage() (uint64, error) {
	lines, err := iox.ReadTextLines("/proc/stat", iox.WithoutBlank())
	if err != nil {
		return 0, err
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			if len(fields) < cpuFields {
				return 0, fmt.Errorf("bad format of cpu stats")
			}

			var totalClockTicks uint64
			for _, i := range fields[1:cpuFields] {
				v, err := parseUint(i)
				if err != nil {
					return 0, err
				}

				totalClockTicks += v
			}

			return (totalClockTicks * uint64(time.Second)) / cpuTicks, nil
		}
	}

	return 0, errors.New("bad stats format")
}

func totalCpuUsage() (usage uint64, err error) {
	var cg *cgroup
	if cg, err = currentCgroup(); err != nil {
		return
	}

	return cg.acctUsageAllCpus()
}
