all:
	@echo "make cnc_win"
	@echo "make cnc_linux"

	@echo "make bot_win"
	@echo "make bot_linux"

cnc_win:
	cd cnc/
	make win_build_cnc

cnc_linux:
	cd cnc/
	make linux_build_cnc

bot_win:
	cd bot/
	make win_build_bot

cnc_linux:
	cd bot/
	make linux_build_bot
