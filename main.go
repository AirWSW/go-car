package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	
	"github.com/AirWSW/gwp"
)

const (
	PIN_PWM_IN1 = gwp.PIN_GPIO_1
	PIN_PWM_IN2 = gwp.PIN_GPIO_4
	PIN_PWM_IN3 = gwp.PIN_GPIO_5
	PIN_PWM_IN4 = gwp.PIN_GPIO_6
	PIN_AVOID_L = gwp.PIN_SPI0_CE1
	PIN_AVOID_R = gwp.PIN_SPI0_CE0
	PIN_TRACK_L = gwp.PIN_GPIO_17
	PIN_TRACK_R = gwp.PIN_GPIO_16
	PIN_MD_MOSI = gwp.PIN_SPI1_MOSI
	PIN_MD_SCLK = gwp.PIN_SPI1_SCLK
)

var mutex = &sync.Mutex{}

func MeasureDistance() (float64, error) {
	gwp.DigitalWrite(PIN_MD_MOSI, gwp.LOW)
	gwp.DelayMicroseconds(2)
	gwp.DigitalWrite(PIN_MD_MOSI, gwp.HIGH)
	gwp.DelayMicroseconds(10)
	gwp.DigitalWrite(PIN_MD_MOSI, gwp.LOW)
	count := 0
	mutex.Lock()
	defer mutex.Unlock()
	for {
		if gwp.DigitalRead(PIN_MD_SCLK) == 1 {
			break
		} else if count < 100000 {
			count++
		} else {
			return 0, errors.New("Failed to read _/ from PIN_MD_SCLK.")
		}
	}
	start := time.Now()
	for {
		if gwp.DigitalRead(PIN_MD_SCLK) == 0 {
			break
		} else if count < 110000 {
			count++
		} else {
			return 0, errors.New("Failed to read \\_ from PIN_MD_SCLK.")
		}
	}
	end := time.Now()
	return float64(end.Sub(start)) / 1000000 * 34000 / 2, nil
}

type Speed struct {
	Pwm1  int
	Pwm2  int
	Pwm3  int
	Pwm4  int
}

func (s *Speed) SetSpeed(p1 int, p2 int, p3 int, p4 int)  {
	s.Pwm1 = p1
	s.Pwm2 = p2
	s.Pwm3 = p3
	s.Pwm4 = p4
}

func main() {
	log.Println("Setup WiringPi...")
	if err := gwp.WiringPiSetup(); err != nil {
		log.Println(err)
	}

	log.Println("Setup PinMode...")
	gwp.PinMode(PIN_PWM_IN1, gwp.OUTPUT)
	gwp.PinMode(PIN_PWM_IN2, gwp.OUTPUT)
	gwp.PinMode(PIN_PWM_IN3, gwp.OUTPUT)
	gwp.PinMode(PIN_PWM_IN4, gwp.OUTPUT)
	gwp.PinMode(PIN_AVOID_L, gwp.INPUT)
	gwp.PinMode(PIN_AVOID_R, gwp.INPUT)
	gwp.PinMode(PIN_TRACK_L, gwp.INPUT)
	gwp.PinMode(PIN_TRACK_R, gwp.INPUT)
	gwp.PinMode(PIN_MD_MOSI, gwp.OUTPUT)
	gwp.PinMode(PIN_MD_SCLK, gwp.INPUT)

	log.Println("Setup InitPwmSpeed...")
	s := Speed{0, 0, 0, 0}

	log.Println("Setup SoftPwmCreate...")
	gwp.SoftPwmCreate(PIN_PWM_IN1, s.Pwm1, 500)
	gwp.SoftPwmCreate(PIN_PWM_IN2, s.Pwm2, 500)
	gwp.SoftPwmCreate(PIN_PWM_IN3, s.Pwm3, 500)
	gwp.SoftPwmCreate(PIN_PWM_IN4, s.Pwm4, 500)
	
	go func() {
		log.Println("Setup MeasureDistance...")
		count := 0
		for {
			dis, err := MeasureDistance()
			if err != nil {
				log.Println(err)
			} else {
				if dis < 10 {
					// s.SetSpeed(0, 250, 0, 250)
					// gwp.Delay(200)
					// s.SetSpeed(250, 0, 0, 250)
					// gwp.Delay(600)
					// log.Println("MD_BACK")
				}
				if dis < 30 {
					// s.SetSpeed(100, 0, 100, 0)
					// log.Println("MD_SLOW")
				}
				// go writer("DISTANCE", []byte(fmt.Sprintf("%0.3f μm", dis)))
				// log.Printf("%0.3f μm", dis)
			}
			if count % 30 == 0 {
				go writer("DISTANCE", []byte(fmt.Sprintf("%0.3f μm", dis)))
			}
			count++
		}
	}()

	go func() {
		log.Println("Setup Avoidance...")
		for {
			AVOID_SL := gwp.DigitalRead(PIN_AVOID_L)
			AVOID_SR := gwp.DigitalRead(PIN_AVOID_R)
			if AVOID_SL == gwp.LOW && AVOID_SR == gwp.LOW {
				// s.SetSpeed(0, 250, 0, 250)
				// gwp.Delay(200)
				// s.SetSpeed(250, 0, 0, 250)
				// gwp.Delay(600)
				// log.Println("AVOID_BACK")
			} else if AVOID_SL == gwp.HIGH && AVOID_SR == gwp.LOW {
				// s.SetSpeed(250, 0, 0, 250)
				// gwp.Delay(600)
				// log.Println("AVOID_RIGHT")
			} else if AVOID_SR == gwp.HIGH && AVOID_SL == gwp.LOW {
				// s.SetSpeed(0, 250, 250, 0)
				// gwp.Delay(600)
				// log.Println("AVOID_LEFT")
			} else {
				// s.SetSpeed(150, 0, 150, 0)
				// log.Println("AVOID_GO")
			}
		}
	}()

	go func() {
		log.Println("Setup Track...")
		for {
			TRACK_SL := gwp.DigitalRead(PIN_TRACK_L)
			TRACK_SR := gwp.DigitalRead(PIN_TRACK_R)
			if TRACK_SL == gwp.LOW && TRACK_SR == gwp.LOW {
				// s.SetSpeed(250, 0, 250, 0)
				// log.Println("TRACK_GO")
			} else if TRACK_SL == gwp.HIGH && TRACK_SR == gwp.LOW {
				// s.SetSpeed(250, 0, 0, 250)
				// log.Println("TRACK_RIGHT")
			} else if TRACK_SR == gwp.HIGH && TRACK_SL == gwp.LOW {
				// s.SetSpeed(0, 250, 250, 0)
				// log.Println("TRACK_LEFT")
			} else {
				// s.SetSpeed(0, 250, 0, 250)
				// log.Println("TRACK_BACK")
			}
		}
	}()

	go func() {
		log.Println("Setup PWM...")
		for {
			/* Left */
			gwp.SoftPwmWrite(PIN_PWM_IN1, s.Pwm1)
			gwp.SoftPwmWrite(PIN_PWM_IN2, s.Pwm2)
			/* Right */
			gwp.SoftPwmWrite(PIN_PWM_IN3, s.Pwm3)
			gwp.SoftPwmWrite(PIN_PWM_IN4, s.Pwm4)
		}
	}()

	go Cam(&s)

	log.Println("Setup WebServer...")
	http.HandleFunc("/subscribe", wsHandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	go echo()

	log.Fatal("ListenAndServe: ", http.ListenAndServe(":8080", nil))

	select {}
}
