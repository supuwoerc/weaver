# wrk压力测试
wrk -t16 -c200 -d15s --latency http://localhost:8804/api/v1/public/captcha/signup