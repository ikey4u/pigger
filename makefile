all:
	@echo "[+] Install to local ..."
	@go install
	@echo "[+] Cross compile for linux, windows and mac ..."
	@packr
	@gox -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="linux windows darwin" -arch="amd64 386"
	@packr clean
	@echo "All is done!"
