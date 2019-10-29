#!/bin/sh

cd /home/pi/mjpg-streamer/mjpg-streamer-experimental/
export LD_LIBRARY_PATH="$(pwd)"
./mjpg_streamer -i "./input_raspicam.so -x 640 -y 480 -fps 30" -o "./output_http.so -p 8081 -w ./www"