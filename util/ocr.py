import ddddocr
import sys
def ocr(file):
    ocr = ddddocr.DdddOcr(beta=True,show_ad=False)
    with open(file, 'rb') as f:
        image = f.read()
    res = ocr.classification(image)
    return res

if __name__ == "__main__":
    file = sys.argv[1]
    print(ocr(file),end="")