// +build dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

/*
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>

#define __USE_GNU       //启用CPU_ZERO等相关的宏
//#define _GNU_SOURCE
#include <sched.h>

int getaffinity(int pid) {
	cpu_set_t mask;
    int i = 0;
	int num = sysconf(_SC_NPROCESSORS_CONF); //获取当前的cpu总数

	CPU_ZERO(&mask);
    if(sched_getaffinity(pid, sizeof(mask), &mask) == -1) {
        return -2;
	}
	for(i = 0; i < num; i++) {
        if(CPU_ISSET(i, &mask)) {
			return i;
		}
	}
	return -1;
}

int setaffinity(int pid, int i) {
	cpu_set_t mask;
	CPU_ZERO(&mask);
	CPU_SET(i, &mask);
    if(sched_setaffinity(pid, sizeof(mask), &mask) == -1) {
        return -2;
	}
	return i;
}
*/
import "C"

import (
	"errors"
	"strconv"
)

var (
	ErrSetAffinity = errors.New("sched_setaffinity fail")
)

func exchangeAffinity(mask, pid int) error {
	if schedGetAffinity(pid) == mask {
		return nil
	}
	// cpu 从0 开始
	mask--
	return schedSetAffinity(pid, mask)
}

func schedGetAffinity(pid int) int {
	return int(C.getaffinity(C.int(pid)))
}

func schedSetAffinity(pid, mask int) error {
	if int(C.setaffinity(C.int(pid), C.int(mask))) != mask {
		return ErrSetAffinity
	}
	return nil
}

func makeAffinities(affinity []string) ([]int, error) {
	affinities := make([]int, len(affinity))
	for i, v := range affinity {
		cpuMask, err := strconv.ParseInt(v, 2, 0)
		if err != nil {
			logf("strconv.ParseInt(%v, 2, 0) error(%v)", v, err)
			return nil, err
		}
		affinityMask := 0
		for cpuMask > 0 {
			cpuMask >>= 1
			affinityMask++
		}
		affinities[i] = affinityMask
	}
	return affinities, nil
}
