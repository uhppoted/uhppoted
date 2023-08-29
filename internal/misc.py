import subprocess


def sublime2():
    return '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'


def say(msg):
    transliterated = msg.replace('uhppoted','u h p p o t e d') \
                        .replace('nodejs','node js') \
                        .replace('codegen', 'code gen') \
                        .replace('Errno','error number') \
                        .replace('exe','e x e') \
                        .replace('unpushed','un pushed') \
                        .replace('cli','c l i') \
                        .replace('github','ggithub')
    subprocess.call(f'say {transliterated}', shell=True)
