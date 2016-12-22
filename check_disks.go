// Checkdisk looks at mounted filesystems, and reports the freespace.
// If the freespace is less than 10%, exit code is non-0.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
)

func main() {
	// capicityThreshold is the percentage when alerts are generated.
	capacityThreshold := 90
	var alert bool
	var output string

	fs := getFs()
	for _, v := range fs {
		var statfs syscall.Statfs_t
		err := syscall.Statfs(v, &statfs)
		if err != nil {
			log.Printf("statfs %s: %s\n", v, err)
		}
		cap := fsCapacity(statfs)
		if cap >= capacityThreshold {
			alert = true
		}
		output = fmt.Sprintf("%s%s: %d%% full\n", output, v, cap)
	}
	if alert == true {
		fmt.Printf("DISK(S) CRITICAL\n%s", output)
		os.Exit(2)
	}
	fmt.Printf("DISK(S) OK\n%s", output)
	os.Exit(0)
}

// getFs returns the names of mounted filesystems.
func getFs() []string {
	f, err := os.Open("/etc/mtab")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	stopwords := []string{"/proc", "/sys", "/dev/pts", "nfs"}
	var fs []string
	for scanner.Scan() {
		// Get the second word in s.Text()
		disk := strings.Fields(scanner.Text())[1]
		var virtualDevice bool
		for _, v := range stopwords {
			if strings.Contains(disk, v) {
				virtualDevice = true
			}
		}
		if !virtualDevice {
			fs = append(fs, disk)
		}
		if err := scanner.Err(); err != nil {
			log.Println("scanner:", err)
		}
	}
	return fs
}

// fsCapacity returns a percentage of how full a filesystem is.
func fsCapacity(fs syscall.Statfs_t) int {
	blocks := float64(fs.Blocks)
	free := float64(fs.Bfree)

	// Capacity is (blocks-free)/blocks * 100
	dif := blocks - free
	quo := dif / blocks
	cap := quo * 100

	return int(cap + 0.5)
}
