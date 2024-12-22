from os import listdir
from PIL import Image
from os.path import isfile, join
mypath = './'
onlyfiles = [f for f in listdir(mypath) if isfile(join(mypath, f))]
images = []
for f in onlyfiles:
    if f.endswith('.png'):
        images.append(f)


for i in images:
    image = Image.open(i)
    pixels = image.load()
    x = image.size[0]
    y = image.size[1]
    for i_x in range(x):
        for i_y in range(y):
            image.putpixel((i_x, 0), (255, 255, 0))
            image.putpixel((0, i_y), (255, 255, 0))
            image.putpixel((i_x, y - 1), (255, 255, 0))
            image.putpixel((x - 1, i_y), (255, 255, 0))

    image.save(i)
