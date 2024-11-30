
.PHONY: default report show before

default: report

report: ./allure-results/.report

show: report
	allure open -p 16258

./allure-results/.report: ./allure-results/*
	-rm -rf ./allure-results/history/
	-cp -rf ./allure-report/history ./allure-results
	-rm -rf ./allure-report
	allure generate
	touch ./allure-results/.report

before:
	rm -rf ./allure-results/*

