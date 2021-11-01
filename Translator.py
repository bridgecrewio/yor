from tabulate import tabulate
from googletrans import Translator


class TranslateClass(object):
    def __init__(self, word, lang):
        self.word = word
        self.lang = lang
        self.Trans = Translator(service_urls=["translate.google.com"])

    def __repr__(self):
        translated = self.Trans.translate(self.word, dest=self.lang).text
        data = [
            ['Language:', "Word/Sentence"],
            ['English', self.word],
            ['Kannada', str(translated)]]
        table = str(tabulate(data, headers="firstrow", tablefmt="grid"))
        return table


if __name__ == '__main__':
    translate = input('Enter Word/Sentence...')
    language = 'kn'  # Translates to Kannada
    # language = 'hi' # Translates to Hindi
    # language = 'it' # Translates to Italian
    print(TranslateClass(translate, language))
