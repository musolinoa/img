MKSHELL=rc

fullsize=`{ls *.full.* >[2]/dev/null}
bigs=`{ls *.full.* | sed 's/\.full\.[A-Z]\+/.big.JPG/'}
mediums=`{ls *.full.* | sed 's/\.full\.[A-Z]\+/.medium.JPG/'}
smalls=`{ls *.full.* | sed 's/\.full\.[A-Z]\+/.small.JPG/'}
thumbs=`{ls *.full.* | sed 's/\.full\.[A-Z]\+/.thumb.JPG/'}
n=`{ls *.full.* | wc -l}
pages=`{seq 1 $n | sed 's/$/.html/'}

all:V: bigs mediums smalls thumbs html montage.jpg
bigs:V: $bigs
mediums:V: $mediums
smalls:V: $smalls
thumbs:V: $thumbs

%.big.JPG:
	convert $stem^.full.* -resize x768 -auto-orient $target

%.medium.JPG: %.big.JPG
	convert $prereq -resize x384 -auto-orient $target

%.small.JPG: %.medium.JPG
	convert $prereq -resize x192 -auto-orient $target

%.thumb.JPG: %.small.JPG
	convert $prereq -resize x96 -auto-orient $target

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
