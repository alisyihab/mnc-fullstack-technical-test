package main

import (
	"fmt"
	"math"
	"mnc-fullstack-technical-test/tahap-1/shared"
)

const OfficeLeave = 14
const NewEmployeeWaitingDays = 180
const MaxConsecutiveLeave = 3

func main() {
	reader := shared.NewInputReader()

	fmt.Print("Input jumlah cuti bersama: ")

	collectiveLeave, ok := reader.ReadInt()
	if !ok {
		fmt.Println("invalid input")
		return
	}

	fmt.Print("Input tanggal join (YYYY-MM-DD): ")

	joinDateStr, ok := reader.ReadLine()
	if !ok {
		fmt.Println("invalid input")
		return
	}

	joinDate, err := shared.ParseDate(joinDateStr)
	if err != nil {
		fmt.Println("invalid join date")
		return
	}

	fmt.Print("Input tanggal rencana cuti (YYYY-MM-DD): ")

	leaveDateStr, ok := reader.ReadLine()
	if !ok {
		fmt.Println("invalid input")
		return
	}

	leaveDate, err := shared.ParseDate(leaveDateStr)
	if err != nil {
		fmt.Println("invalid leave date")
		return
	}

	fmt.Print("Input durasi cuti: ")

	duration, ok := reader.ReadInt()
	if !ok {
		fmt.Println("invalid input")
		return
	}

	if duration > MaxConsecutiveLeave {
		fmt.Println("False")
		fmt.Println("Alasan: Cuti pribadi max 3 hari berturutan")
		return
	}

	personalLeave := OfficeLeave - collectiveLeave

	eligibleDate := shared.AddDays(joinDate, NewEmployeeWaitingDays)

	if leaveDate.Before(eligibleDate) {
		fmt.Println("False")
		fmt.Println("Alasan: Karena belum 180 hari sejak tanggal join karyawan")
		return
	}

	// Tahun pertama
	if joinDate.Year() == leaveDate.Year() {

		remainingDays := shared.DaysUntilEndOfYear(eligibleDate)

		quota := int(math.Floor(
			(float64(remainingDays) / 365.0) * float64(personalLeave),
		))

		if duration > quota {
			fmt.Println("False")
			fmt.Printf("Alasan: Karena hanya boleh mengambil %d hari cuti\n", quota)
			return
		}
	}

	fmt.Println("True")
}
