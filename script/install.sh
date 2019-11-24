mkdir -p $HOME/ticktak
cp -r ../* $HOME/ticktak
systemctl --user enable $HOME/ticktak/bin/ticktak_server
systemctl --user start $HOME/ticktak/bin/ticktak_server