#!/bin/sh
for i in "WenQuanYiZenHei-01.ufo3" "WenQuanYiZenHeiMono-02.ufo3" "WenQuanYiZenHeiSharp-03.ufo3"; do
  /usr/bin/fontmake -u $i --validate-ufo --overlaps-backend pathops -a --verbose DEBUG -o otf-cff2 ttf
done
