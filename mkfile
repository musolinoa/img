MKSHELL=/bin/rc

years=`{ls | grep '^[0-9][0-9][0-9][0-9]$'}

dirs=\
	moto\
	guns\
	sssc\
	group-shoot\
	ttc2022\

fullsize=`{find $years -type f -name '*.full.*'}

montages=`{echo $years/^montage.jpg}

index.html:D: $montages ./mkindex.rc ./mkfile
	./mkindex.rc $years >index.html

%/montage.jpg: subdirs

subdirs:V:
	for(d in $years)@{
		cd $d
		mk -f ../year.mk
	}

html clean nuke:V:
	for(d in $years)@{
		cd $d
		mk -f ../year.mk $target
	}
