# 다음과 같이 명령어 실행해보세요
go run receive_logs_direct.go warning error > logs_from_rabbit.log
go run receive_logs_direct.go info warning error
go run emit_log_direct.go error "Run. Run. Or it will explode."