
DIR_TEST_REPORT := allure-report/

.PHONY: default report show before

default: report

report: $(DIR_TEST_REPORT)/.report

show: report
	cd $(DIR_TEST_REPORT) && allure open -p 16258

$(DIR_TEST_REPORT)/.report: $(DIR_TEST_REPORT)/allure-results/*
	mkdir -p $(DIR_TEST_REPORT)
	-rm -rf $(DIR_TEST_REPORT)/allure-results/history/
	-cp -rf $(DIR_TEST_REPORT)/allure-report/history $(DIR_TEST_REPORT)/allure-results
	-rm -rf $(DIR_TEST_REPORT)/allure-report
	cd $(DIR_TEST_REPORT) && allure generate
	touch $(DIR_TEST_REPORT)/.report

before:
	rm -rf $(DIR_TEST_REPORT)/.report
	rm -rf $(DIR_TEST_REPORT)/allure-results

