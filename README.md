  # file naming conventions
  
- Allowed Metacharacters are (**double colon** `::`,   **underscore** `_` and **period** `.`) . Everything else must not be used in file name such as:
	  -	`'`, `?`, `#`,`"` , `...`
- If using a book name from the bible that has 2 words like "1 Timothy", truncate it removing the space.
	- for example: `::book.1timothy::`
- filename will be in sections.
- sections will be delimited by the double colon `::`
- spaces will be substituted by underscore `_`
	- ex: `this_is_four_words`

- filenames will have a **field** **value** structure like `name.john_verstegen`
- filename `field value`'s will be seperated by the double colon like:  `name.john_verstegen::date.20220508::.mp4`
- filenames will always end with filetype postfix like:  `name.john_verstegen::.mp4` emphasis on **.mp4**
Example filename: **topic.right_division::date.20220507::speaker.john_verstegen::title.where_should_i_start::.mp4**

### braindump

- old filename: `20210822-JAMES_SERIES-PART_75-JOHN_VERSTEGEN-JAMES_CH_5_V_10_AND_11-JOB_CH_38_V_1_TO_12-CONFINING_THE_REBELLION.mp3`
- new filename: `date::20210822::subject.james::part.75::speaker.john_verstegen::book.(james,job)::chapter.(5,38)::verse.(5,1-12)::title.confining_the_rebellion.mp3`

-- api response -- 
```json
[
    {"date": 20210822,
     "subject": "james",
     "series": "james",
     "part": 75,
     "speaker": "john verstegen",
     "book":[
         {
             "name": "james",
             "start_chapter": 5
             "end_chapter": null,
             "start_verse": 5,
             "end_verse": null,
         },
         {
             "name": "job",
             "start_chapter": 38,
             "end_chapter": null,
             "start_verse": 1,
             "end_verse": 12,
         }
     ]
     "title": "confining the rebellion",
     "type": "mp3"
    }
]
```
