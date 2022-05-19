nothing:
	@echo "make win_build_cnc"
	@echo "make linux_build_cnc"

linux_build_bot: linux_clean
	# cd bot/
	go build ./src
	chmod +x ./src
	clear
	./src

win_build_bot: win_clean
	# cd bot/
	go build ./src
	cls
	./src.exe

win_clean:
	del *.exe
linux_clean:
	rm ./src
