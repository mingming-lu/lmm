import argparse
import os
import sys

from PIL import Image


def resize_single_file(dst: str, src: str, width: int):
    with Image.open(src) as img:
        if img.size[0] < width:
            print(f'skip resizing image {src} since width ({img.size[0]}) is smaller than {width}')
            img.save(dst)
            return
        ratio = width / img.size[0]
        height = int(img.size[1] * ratio)
        img.thumbnail((width, height), Image.ANTIALIAS)
        img.save(dst)
        print(f'resized {src} wrote to {dst}')


def resize_all_files_from_dir(dst_dir: str, src_dir: str, width: int):
    os.makedirs(os.path.dirname(dst_dir), mode=0o644, exist_ok=True)
    for file in os.listdir(src_dir):
        if file.endswith('.jpeg') or file.endswith('.jpg') or file.endswith('.png'):
            dst = os.path.join(dst_dir, file)
            src = os.path.join(src_dir, file)
            resize_single_file(dst, src, width)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--src', help='image src')
    parser.add_argument('--dst', help='image dst')
    parser.add_argument('--src_dir', help='image src dir')
    parser.add_argument('--dst_dir', help='image dst dir')
    parser.add_argument('--width', type=int, required=True, help='specify width to resize with the same ratio')

    args = parser.parse_args()

    if args.dst is not None and args.src is not None:
        resize_single_file(args.dst, args.src, args.width)

    if args.dst_dir is not None and args.src_dir is not None:
        resize_all_files_from_dir(args.dst_dir, args.src_dir, args.width)
