package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"time"

	"gocv.io/x/gocv"
)

func Cam(s *Speed) {
	// set to use a video capture device 0
	deviceID := 0

	// http://0.0.0.0:8081/?action=snapshot&n=0
	count := 0

	// prepare image matrix
	mat := gocv.NewMat()
	defer mat.Close()

	// prepare gray image matrix
	gray := gocv.NewMat()
	binary := gocv.NewMat()
	defer gray.Close()
	defer binary.Close()

	log.Printf("start reading camera device: %v\n", deviceID)
	
	for {
		t1 := time.Now() // get current time
		url := fmt.Sprintf("http://0.0.0.0:8081/?action=snapshot&n=%d", count)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body) // 2~7ms
		if err!=nil {
			log.Println(err)
		}
		mat, err := gocv.IMDecode(buf, gocv.IMReadGrayScale) // 40~42ms
		if err != nil {
			log.Println(err)
		}

		if mat.Empty() {
			continue
		}
		
		log.Println("start to detect lines...")
		gocv.Threshold(mat, &binary, 50, 255, gocv.ThresholdBinaryInv)
		matBytes := binary.ToBytes()

		sumA := 0
		sumB := 0
		sumC := 0
		countA := 0
		countB := 0
		countC := 0
		for i := 0; i < 640; i++ {
			if matBytes[640 * 150 + i] == 255 {
				sumA = sumA + i
				countA++
			}
			if matBytes[640 * 250 + i] == 255 {
				sumB = sumB + i
				countB++
			}
			if matBytes[640 * 250 + i] == 255 {
				sumC = sumC + i
				countC++
			}
		}
		if countA == 0 {
			countA = 1
		}
		if countB == 0 {
			countB = 1
		}
		if countC == 0 {
			countC = 1
		}
		if sumA/countA > sumB/countB + 10 {
			s.SetSpeed(200, 0, 100, 0)
		} else if sumA/countA < sumB/countB - 10 {
			s.SetSpeed(100, 0, 200, 0)
		} else {
			s.SetSpeed(120, 0, 120, 0)
		}
		// log.Println(sumA/countA, sumB/countB, sumC/countC)
		
		contours := gocv.FindContours(binary, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		approxs := contours
		for n, curve := range contours {
			approxs[n] = gocv.ApproxPolyDP(curve, 3, true)
		}

		elapsed := time.Since(t1)
		
		log.Println("Got: ", url, ", elapsed: ", elapsed, fmt.Sprintf(", %d end", count))
		go writer("CONTOURS", []byte(fmt.Sprintf("[[(%d,150) (%d,250) (%d,250)]]", sumA/countA, sumB/countB, sumC/countC)))
		go writer("CONTOURS", []byte(fmt.Sprintln(approxs)))
		go writer("CAMINFO", []byte(fmt.Sprintln("Got: ", url, ", elapsed: ", elapsed, fmt.Sprintf(", %d end", count))))
		count++
	}
}