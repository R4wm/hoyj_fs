#!/usr/bin/env python3
import json
import re

# create static speaker list
speakers = ['alex_kurz', # acts-series-alex-kurz
            'barry_curtis', 'bryan_ross', 'dave_stout',
            'david_busch', 'david_stout', 'elbert_ray', 'eric_neumann',
            'henry_meneses',
            'garrett_kaylor', 'gene_fuqua', 'jeff_halcomb', 'jim_lawrence',
            'joe_purczynski', 'john_versegen', 'john_verstegen',
            'keith_baxter', 'morris_chestnut', 'nate_cody', 'paul_zoerb',
            'rich_davis', 'rich_jordan', 'richard_jordan', 'rick_jordan',
            'rick_jordan', 'rom_knight', 'ron_knight', 'ted_follows',
            'ted_marolf', 'tom_meier',
            'willard' # acts-series--bro-willard
            ]

bookList = ['genesis', 'exodus', 'leviticus', 'numbers', 'deuteronomy', 'joshua', 'judges', 'ruth', '1samuel', '2_samuel', '1_kings', '2_kings', '1_chronicles', '2_chronicles', 'ezra', 'nehemiah', 'esther', 'job', 'psalms', 'proverbs', 'ecclesiastes', 'song of solomon', 'isaiah', 'jeremiah', 'lamentations', 'ezekiel', 'daniel', 'hosea', 'joel', 'amos', 'obadiah', 'jonah', 'micah', 'nahum', 'habakkuk', 'zephaniah', 'haggai', 'zechariah', 'malachi', 'matthew', 'mark', 'luke', 'john', 'acts', 'romans', '1_corinthians', '2_corinthians', 'galatians', 'ephesians', 'philippians', 'colossians', '1_thessalonians', '2_thessalonians', '1_timothy', '2_timothy', 'titus', 'philemon', 'hebrews', 'james', '1_peter', '2_peter', '1_john', '2_john', '3_john', 'jude', 'revelation']

# Make array of MD5sum <-> filename
MD5TOFILE = []  # contains tuples(md5, filename)
md5list = open("list.md5", "r")
for line in md5list:
    # print(line)
    md5sum, md5filename = line.split('  ')
    MD5TOFILE.append((md5sum, md5filename.strip('\n').lower()))
md5list.close()


# Get all filenames topic/filename
MEDIALIST = []
medialist = open("media.txt", 'r')
for line in medialist:
    MEDIALIST.append(line.strip('\n'))
medialist.close()


# read config
with open("labeling_conf.json") as r:
    APPCONFIG = r.read()
    APPCONFIG = json.loads(APPCONFIG)


class Message:
    """the reason we are here"""
    def __init__(self, filename=''):
        self.originalref = filename
        filename = filename.lower()
        self.filename = filename
        self.speaker = ''
        self.year = None
        self.month = None
        self.day = None
        self.topic = ''
        self.md5sum = ''
        self.book = ''
        self.chapter = None
        self.part = None
        self.verseStart = None
        self.verseEnd = None
        self.outlinefile = ''

    def setSpeaker(self):
        if self.speaker != '':
            return
        for i in speakers:
            if i in self.filename.lower():
                self.speaker = i.replace('_', ' ')

    def setTopic(self):
        if '/' not in self.filename:
            print("skipping, no slash found in filename")
            return
        topic = self.filename.split('/')[0]
        if topic == 'acts-series--bro-willard':
            self.topic = 'acts-series'
            self.speaker = 'willard'
            self.book = 'acts'
            return
        if topic == 'acts-series-alex-kurz':
            self.topic = 'acts-series'
            self.speaker = 'alex kurz'
            return
        if 'single-message' in topic:
            self.topic = 'single-message'
            return
        endswithyear = re.match(r'(.*)(-\d{4})', topic)
        if endswithyear:
            self.topic = topic.replace(endswithyear.group(2), '')
            self.year = int(endswithyear.group(2).replace('-', ''))
            return
        if topic == 'romans-chapter-15-series':
            self.topic = topic
            self.chapter = 15
            return
        startswithyear = re.match(r'(\d{4})(-)(.*)', topic)
        if startswithyear:
            if self.topic == '':
                self.topic = startswithyear.group(3)
                return
        # default to whatever is there
        self.topic = topic

    def setDate(self):
        foundDate = re.match(r'(.*)(\d{8})', self.filename)
        if foundDate:
            self.year = int(foundDate.group(2)[0:4])
            self.month = int(foundDate.group(2)[4:6])
            self.day = int(foundDate.group(2)[6:])

    def setMD5Sum(self):
        for i in MD5TOFILE:
            checkAgainst = self.filename.split('/')[1]
            if i[1] == checkAgainst:
                self.md5sum = i[0]
                return

    def setChapter(self):
        checkAgainst = self.filename.split('/')[1]
        foundChapter = re.match(r'(.*)(_ch_)(\d{1,3})', checkAgainst)
        if foundChapter:
            if self.chapter is None:
                self.chapter = int(foundChapter.group(3))
                return
        cornerCase1 = re.match(r'(.*-_)([a-z]+)(_ch\._)(\d{1,3})', checkAgainst)
        if cornerCase1:
            if self.chapter is None:
                self.chapter = int(cornerCase1.group(4))


    def setBook(self):
        checkAgainst = self.filename.split('/')[1]
        foundBook = re.match(r'(.*[-_])(.*)(_ch_\d{1,3})', checkAgainst)
        if foundBook:
            if self.book == '':
                self.book = foundBook.group(2)
                return
        cornerCase1 = re.match(r'(.*-_)([a-z]+)(_ch\._)(\d{1,3})', checkAgainst)
        if cornerCase1:
            if self.book == '':
                self.book = cornerCase1.group(2)
                return
        # last resorts
        if self.book == '':
            for i in ['2_corinthians', '1_corinthians', '1_timothy',
                      '2_timothy', '1_samuel', '2_samuel', '1_kings',
                      '2_kings', '1_chronicals', '2_chronicals',
                      '1_thessalonians', '2_thessalonians', '1_peter', '2_peter',
                      '1_john', '2_john', '3_john']:
                if i in self.filename:
                    self.book = i.replace('_', ' ')
                    return

    def setVerses(self):
        checkAgainst = self.filename.split('/')[1]
        # luke_ch_15_v_11_to_32.mp3
        foundVerses = re.match(r'(.*)(_ch_\d{1,3})(_)(v_)(\d{1,3})(_\w{1,3}_)(\d{1,3})', checkAgainst)
        if foundVerses:
            if self.verseStart is None and self.verseEnd is None:
                self.verseStart = int(foundVerses.group(5))
                self.verseEnd = int(foundVerses.group(7))
                return
        # handle single verse ex: john_verstegen-philippians_ch_3_v_17-who_to_follow_and_who_to_avoid.mp3
        foundVerse = re.match(r'(.*)(_ch_\d{1,3})(_)(v_)(\d{1,3})([-_].*)', checkAgainst)
        if foundVerse:
            if self.verseStart is None:
                self.verseStart = int(foundVerse.group(5))
                return
        # cornerCase1 = re.match(r'(.*-_)([a-z]+)(_ch\._)(\d{1,3})', checkAgainst)
        # if cornerCase1:
        #     import pdb;pdb.set_trace()


    def setPart(self):
        checkAgainst = self.filename.split('/')[1]
        foundPart = re.match(r'(.*[_-]part_)(\d{1,3})', checkAgainst)
        if foundPart:
            if self.part is None:
                self.part = int(foundPart.group(2))
                return
        cornerCase1 = re.match(r'(.*)(_pt_)(\d{1,3})([-_])', checkAgainst)
        if cornerCase1:
            if self.part is None:
                self.part = int(cornerCase1.group(3))
                return
        # default
        self.part = None

    def lastMinuteChanges(self):
        if self.book == 'timothy':
            if '1_timothy' in self.filename:
                self.book = '1 timothy'
            elif '2_timothy' in self.filename:
                self.book = '2 timothy'
        if self.book == 'series':
            if '1_corinthians' in self.filename:
                self.book = '1 corinthians'
        if self.book == 'series':
            for i in bookList:
                if i in self.filename.replace('john_verstegen', ''):
                    self.book = i.replace('_', ' ')
                    break

    def shortenfilename(self):
        if '/' in self.filename:
            self.filename = self.filename.split('/')[1]


# function to identify year, return Y-M-D

# function to get topic from .split('/')[0]

# function to get the md5sum of the file


# iterate the media.txt

def main():

    # # test topic starting with year
    # s = Message(filename='2011-rightly-dividing-training-series/20101212-SUN-1000-RIGHTLY_DIVIDING_TRAINING_SERIES-GODS_PROPHECIED_PLAN_TO_RECONCILE_THE_EARTH-JOHN_VERSTEGEN.mp3')
    # s.setTopic()
    # s.setSpeaker()
    # s.setMD5Sum()
    # s.setBook()
    # s.setChapter()
    # s.setVerses()
    # s.setPart()
    # print(s.__dict__)

    # t = Message(filename='acts-series--bro-willard/Acts_21_v_24_to_40_ch_22_v_1_to_15_Jun_2_2011.mp3')
    # t.setTopic()
    # t.setSpeaker()
    # t.setMD5Sum()
    # t.setBook()
    # t.setChapter()
    # t.setVerses()
    # t.setPart()
    # print(t.__dict__)

    # # testing topic with year end
    # print("doing u")
    # u = Message(filename='summer-conference-2007/20070801-WED-1900-SUMMER__CONFERENCE-FINDING_REST_IN_A_WORLD_OF_UNREST-RICHARD_JORDAN.mp3')
    # u.setTopic()
    # u.setSpeaker()
    # u.setMD5Sum()
    # u.setBook()
    # u.setChapter()
    # u.setVerses()
    # u.setPart()
    # print(u.__dict__)

    # # testing topic with year middle
    # print("doing v")
    # v = Message(filename='philippians-2019-series/20210314-PHILIPPIANS_SERIES-PART_97-JOHN_VERSTEGEN-PHILIPPIANS_CH_3_V_17-WHO_TO_FOLLOW_AND_WHO_TO_AVOID.mp3')
    # v.setTopic()
    # v.setSpeaker()
    # v.setDate()
    # v.setMD5Sum()
    # v.setBook()
    # v.setChapter()
    # v.setVerses()
    # v.setPart()
    # print(v.__dict__)

    # # testing setBook
    # print("doing w")
    # w = Message(filename='single-message-2/20210620-SUN-1100-THE_LOVE_OF_OUR_FATHER-GARRETT_KAYLOR_LUKE_CH_15_V_11_TO_32.mp3')
    # w.setTopic()
    # w.setSpeaker()
    # w.setDate()
    # w.setMD5Sum()
    # w.setChapter()
    # w.setBook()
    # w.setVerses()
    # w.setPart()
    # print(w.__dict__)

    # # testing setBook
    # print("doing x")
    # x = Message(filename='will-of-god-series/20160518-WED-1900-WILL_OF_GOD_SERIES-PART_4-JOHN_VERSTEGEN_–_1_TIMOTHY_CH_2_V_3_TO_6__WHO_WILL_HAVE_ALL_MEN_TO_BE_SAVED_AND_TO_COME_UNTO_THE_KNOWLEDGE_OF_THE_TRUTH.mp3')
    # x.setTopic()
    # x.setSpeaker()
    # x.setDate()
    # x.setMD5Sum()
    # x.setChapter()
    # x.setBook()
    # x.setVerses()
    # x.setPart()
    # print(x.__dict__)

    # # testing verses with "AND"
    # print("doing y")
    # y = Message(filename='2021-james-series/20210110-JAMES_SERIES-PART_49_JOHN_VERSTEGEN-JAMES_CH_3_V_17_AND_18-WISDOM_FROM_ABOVE.mp3')
    # y.setTopic()
    # y.setSpeaker()
    # y.setDate()
    # y.setMD5Sum()
    # y.setChapter()
    # y.setBook()
    # y.setVerses()
    # y.setPart()
    # print(y.__dict__)

    # # cornercase 1
    # print("doing z")
    # z = Message(filename='winter-conference-2018/20180217-SAT-0900-WINTER_CONFERENCE_2018-RICK_JORDAN_-_Psalm_Ch._16_Luke_Ch._24_V._5_and_6._The_Message_Of_An_Empty_Tomb.mp3')
    # z.setTopic()
    # z.setSpeaker()
    # z.setDate()
    # z.setMD5Sum()
    # z.setChapter()
    # z.setBook()
    # z.setVerses()
    # z.setPart()
    # print(z.__dict__)

    result = []
    for mediaFile in MEDIALIST:
        a = Message(filename=mediaFile)
        a.setTopic()
        a.setSpeaker()
        a.setDate()
        a.setMD5Sum()
        a.setChapter()
        a.setBook()
        a.setVerses()
        a.setPart()
        a.lastMinuteChanges()
        a.shortenfilename()
        print(json.dumps(a.__dict__, indent=4))
        result.append(a.__dict__)

    # write the result to file
    with open('/tmp/result.json', 'w') as w:
        w.write(json.dumps(result, indent=4))


if __name__ == '__main__':
    main()
