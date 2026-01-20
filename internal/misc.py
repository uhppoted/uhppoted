import subprocess


def editor():
    return '"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"'


def say(msg):
    transliterated = msg.replace('uhppoted','u h p p o t e d') \
                        .replace('uhppote','u h p p o t e') \
                        .replace('nodejs','node js') \
                        .replace('codegen', 'code gen') \
                        .replace('Errno','error number') \
                        .replace('exe','e x e') \
                        .replace('unpushed','un pushed') \
                        .replace('cli','c l i') \
                        .replace('github','ggithub') \
                        .replace('.10','.ten')
    subprocess.call(f'say {transliterated}', shell=True)
