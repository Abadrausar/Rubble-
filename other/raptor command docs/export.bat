
@del %1.txt
:: Generate docs only for items beginning with "Command"
@godoc dctech/raptor/commands/%1 Command* >> %1.txt
