MKSHELL=rc

months=`{ls | grep '^[0-1][0-9]$' >[2]/dev/null}
montages=`{ls -d $months | sed 's,$,/montage.jpg,'}

all:V: $montages index.html montage.jpg subdirs

[0-9]+/montage\.jpg:RQ:
	for(d in $months)@{
		if(test -d $d){
			cd $d
			mk -f ../../album.mk
		}
		if not
			status=()
	}

subdirs:V:
	for(d in $months)@{
		cd $d
		mk -f ../../album.mk
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
