mkdir -p $HOME/ticktak
cp -r ../* $HOME/ticktak
systemctl --user enable $HOME/ticktak/script/ticktak.service
systemctl --user start $HOME/ticktak/script/ticktak.service