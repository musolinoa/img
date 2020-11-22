MKSHELL=rc

fullsize=`{ls *.full.JPG >[2]/dev/null}
bigs=`{ls *.full.JPG | sed 's/\.full\.JPG/.big.JPG/'}
mediums=`{ls *.full.JPG | sed 's/\.full\.JPG/.medium.JPG/'}
smalls=`{ls *.full.JPG | sed 's/\.full\.JPG/.small.JPG/'}
thumbs=`{ls *.full.JPG | sed 's/\.full\.JPG/.thumb.JPG/'}
n=`{ls *.full.JPG | wc -l}
pages=`{seq 1 $n | sed 's/$/.html/'}

cflags=-auto-orient

all:V: bigs mediums smalls thumbs html montage.jpg
bigs:V: $bigs
mediums:V: $mediums
smalls:V: $smalls
thumbs:V: $thumbs

%.big.JPG: %.full.JPG
	convert $prereq -resize 1024 -auto-orient $target

%.medium.JPG: %.big.JPG
	convert $prereq -resize 512 -auto-orient $target

%.small.JPG: %.medium.JPG
	convert $prereq -resize 256 -auto-orient $target

%.thumb.JPG: %.small.JPG
	convert $prereq -resize 128 -auto-orient $target

montage.jpg: smalls
	$HOME/img/mkmontage.rc

html:V: index.html $pages

index.html: $thumbs $HOME/img/mkalbumindex.rc
	$HOME/img/mkalbumindex.rc >$target

$pages: $HOME/img/mkpages.rc
	$HOME/img/mkpages.rc

clean.html:V:
	rm -f *.html

clean.bigs:V:
	rm -f *.big.JPG

clean.mediums:V:
	rm -f *.medium.JPG

clean.smalls:V:
	rm -f *.small.JPG

clean.thumbs:V:
	rm -f *.thumb.JPG

clean.imgs:V: clean.bigs clean.mediums clean.smalls clean.thumbs
	rm -f montage.jpg

clean:V: clean.html clean.imgs
