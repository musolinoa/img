MKSHELL=rc

months=\
	01\
	02\
	03\
	04\
	05\
	06\
	07\
	08\
	09\
	10\
	11\
	12\

months=`{ls -d $months >[2]/dev/null}
montages=`{ls -d $months | sed 's,$,/montage.jpg,'}

all:V: $montages index.html montage.jpg

[0-9]+/montage\.jpg:RQ:
	for(d in $months)@{
		if(test -d $d){
			cd $d
			mk -f ../../album.mk
		}
		if not
			status=()
	}

html:V: index.html
	for(d in $months)@{
		cd $d
		mk -f ../../album.mk $target
	}

index.html: $HOME/img/mkyearidx.rc $montages
	$HOME/img/mkyearidx.rc >index.html

montage.jpg: $HOME/img/mkmontage.rc $montages
	$HOME/img/mkmontage.rc

clean:V:
	rm *.html
	for(d in $months)@{
		cd $d
		mk -f ../../album.mk $target
	}

nuke:V: clean
	for(d in $months)@{
		cd $d
		mk -f ../../album.mk $target
	}
