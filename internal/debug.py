#!python3

import re
import signal
import traceback
import subprocess

from misc import say


def main():
    print()
    print("*** debug.py")
    print()

    try:
        with open('./uhppote-core/CHANGELOG.md', 'r', encoding='utf-8') as f:
            changelog = f.read()
            matches = list(re.finditer(r'(?m)^## \[[^\]]+\].*?(?=^## \[|\Z)', changelog,
                                       flags=re.DOTALL | re.MULTILINE))

            if matches:
                release_notes = matches[0].group().strip() + '\n'
                subprocess.run("pbcopy", text=True, input=release_notes)
            else:
                subprocess.run("pbcopy", text=True, input='')

    except BaseException as x:
        print(traceback.format_exc())

        msg = f'{x}'
        msg = msg.replace('uhppoted-','')                         \
                 .replace('uhppote-','')                          \
                 .replace('uhppoted','umbrella project')          \
                 .replace('cli','[[char LTRL]]cli[[char NORM]]')  \
                 .replace('git','[[inpt PHON]]git[[input TEXT]]') \
                 .replace('codegen','code gen')

        print()
        print(f'*** ERROR  {x}')
        print()

        say('ERROR')
        say(msg)


if __name__ == '__main__':
    for sig in ['SIGTERM', 'SIGHUP', 'SIGINT']:
        signal.signal(getattr(signal, sig), quit)

    main()
