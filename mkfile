MKSHELL=/bin/rc

years=\
	2008\
	2009\
	2013\
	2014\
	2015\
	2018\
	2019\
	2020\

dirs=\
	moto\
	guns\
	sssc\
	group-shoot\

fullsize=`{find $years $dirs -type f -name '*.full.*'}

montages=`{echo $years/^montage.jpg $dirs/^montage.jpg}

index.html:D: $montages ./mkindex.rc ./mkfile
	./mkindex.rc $years $dirs >index.html

%/montage.jpg: subdirs

subdirs:V:
	for(d in $years)@{
		cd $d
		mk -f ../year.mk
	}
	for(d in $dirs)@{
		cd $d
		mk -f ../album.mk
	}

html clean nuke:V:
	for(d in $years)@{
		cd $d
		mk -f ../year.mk $target
	}
	for(d in $dirs)@{
		cd $d
		mk -f ../album.mk $target
	}
