# -*- coding: utf-8 -*- 
#!/usr/bin/env python
import matplotlib.pyplot as plt
import logging as log
import pandas as pd
import numpy as np
import click
import os
import re

log.basicConfig(format='[%(asctime)s][%(levelname)s] %(message)s', level=log.INFO)

@click.command()
@click.option("-f", "--filename", help="log file name", type=str, required=True)
def main(filename):
	log.info("analyze log file: {}".format(filename))
	pattern = re.compile(r".*\[interval\]:\s(-?\d+)\(ms\)$", re.I)
	l = []
	with open(filename, 'r') as reader:
		for line in reader.readlines():
			line = line.rstrip()
			m = pattern.match(line)
			if m == None:
				continue
			if int(m.group(1)) > 0:
				l.append(int(m.group(1)))
	log.info("len of interval: {}".format(len(l)))
	df = pd.DataFrame(np.array(l))
	print(df.describe())
	# df.loc[df[0]>500] = 500
	# df.plot()
	# plt.show()
	pass

#
# 主程序
if __name__ == '__main__':
	main()
