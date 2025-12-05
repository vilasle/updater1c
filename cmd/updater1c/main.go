package main

/*
	need to have database with actual state of infobases, link of connection, current version

	need way for getting current version, may be use oscript or efp?

	check updates on 1c update-service, download and

	lock new session
	lock background tasks
	delete current session
	create access word

	create dump
	delete all of extensions
	load cfu on infobase
	update
	run 1c and mark that updates was got legal
	run efp for running postponed updates
	if all if fine
	update to next version

	in the finish download and load extensions
	update information in table with infobases

	if some went wrong, restore dump

	create report and send message on telegram

*/
