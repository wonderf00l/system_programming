package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

// BenchmarkSyscallRead benchmarks reading the file portion by portion using syscall.Read
func BenchmarkSyscallRead(b *testing.B) {
	// Create a temporary file with 100MB of data for testing
	filePath := "/tmp/testfile"
	file, err := os.Create(filePath)
	if err != nil {
		b.Fatalf("unable to create test file: %v", err)
	}
	defer file.Close()

	// Fill the file with some data (100MB)
	data := make([]byte, 100*1024*1024) // 100MB of data
	_, err = file.Write(data)
	if err != nil {
		b.Fatalf("unable to write to file: %v", err)
	}

	// Open the file for reading using syscall
	fd, err := syscall.Open(filePath, syscall.O_RDONLY, 0)
	if err != nil {
		b.Fatalf("unable to open file with syscall: %v", err)
	}
	defer syscall.Close(fd)

	// Get file size using fstat
	var stat syscall.Stat_t
	err = syscall.Fstat(fd, &stat)
	if err != nil {
		b.Fatalf("unable to get file stat: %v", err)
	}
	fileSize := stat.Size

	// Benchmark the syscall Read method
	buffer := make([]byte, 64*1024) // 64KB buffer

	b.ResetTimer()
	b.ReportAllocs()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var totalRead int64
		for totalRead < fileSize {
			// Read the file using syscall.Read
			n, err := syscall.Read(fd, buffer)
			if err != nil && err.Error() != "EOF" {
				b.Fatalf("read error: %v", err)
			}
			if n == 0 {
				// EOF reached
				break
			}
			totalRead += int64(n)
		}
	}
	b.StopTimer()
}

// BenchmarkMmap benchmarks memory-mapping file portions using mmap syscall
func BenchmarkMmap(b *testing.B) {
	// Create a temporary file with 100MB of data for testing
	filePath := "/tmp/testfile"
	file, err := os.Create(filePath)
	if err != nil {
		b.Fatalf("unable to create test file: %v", err)
	}
	defer file.Close()

	// Fill the file with some data (100MB)
	data := make([]byte, 100*1024*1024) // 100MB of data
	_, err = file.Write(data)
	if err != nil {
		b.Fatalf("unable to write to file: %v", err)
	}

	// Open the file for reading
	file, err = os.Open(filePath)
	if err != nil {
		b.Fatalf("unable to open file: %v", err)
	}
	defer file.Close()

	// Get file size
	info, err := file.Stat()
	if err != nil {
		b.Fatalf("unable to get file stat: %v", err)
	}
	fileSize := info.Size()

	// Define the chunk size to mmap portion-by-portion
	chunkSize := 64 * 1024 // 64KB

	// Benchmark the mmap method
	b.ResetTimer()
	b.ReportAllocs()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var totalRead int64
		for totalRead < fileSize {
			// Memory-map a portion of the file
			offset := totalRead
			remaining := fileSize - totalRead
			// Map only the remaining portion, limited by chunkSize
			if remaining < int64(chunkSize) {
				chunkSize = int(remaining)
			}

			data, err := mmap(file, int(offset), chunkSize)
			if err != nil {
				b.Fatalf("mmap failed: %v", err)
			}

			// Read the mapped data (this is a no-op, just accessing the memory)
			_ = data

			// Update totalRead
			totalRead += int64(len(data))

			// Unmap the portion (this is implicit in Go, but we'll simulate the idea)
			// In Go, we don't explicitly need to unmap, but we do this manually to match the logic
			// on systems that require explicit unmapping
			// syscall.Munmap(data) // This is typically unnecessary for this Go-based approach.
		}
	}
	b.StopTimer()
}

// mmap memory-maps a portion of the file to memory and returns the byte slice
func mmap(file *os.File, offset int, size int) ([]byte, error) {
	fd := file.Fd()
	data, err := syscall.Mmap(int(fd), int64(offset), size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %w", err)
	}
	return data, nil
}
