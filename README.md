# CLI TODO app
### Create TODO's without leaving the command line!

This is a [Golang](https://golang.org/) app created by using [Cobra](https://github.com/spf13/cobra) package for CLI and [sqlite](https://github.com/mattn/go-sqlite3) package for connecting to DB


## Usage

- Initiate the app by running
    ```bash
    todo init
    ```
    > Note: To be run only once
- Create a note
    ```bash
    todo notes create -m "An important thing"
    ```

  Give it an expected date to finish it by
    ```bash
    todo notes create -e tomorrow -m "Time sensitive note to do by tomorrow"
    ```

  Use your favourite editor by setting `EDITOR` env variable
    ```bash
    export EDITOR=/usr/bin/vim
    todo notes create -e +10d
    ```
  This would open up the editor defined in the `EDITOR` variable

- List notes
    ```bash
    todo notes
    ```
  This would list all notes under the sun (Just kidding, only the ones created)

  List notes by expected notes
    ```bash
    todo notes -A yesterday -B +15d
    ```
    This would get notes after yesterday and before 15 days

  Get the created time too
    ```bash
    todo notes --created
    ```
- Dates: To make things a bit easier, certain formats of dates can be passed
  - The actual date like `YYYY-MM-DD`. eg `2020-01-02`
  - English: `today`, `tomorrow`, `yesterday`, `day-after` or `day_after`
  - Counting days, week, month or year:
    - `+5d`: Plus five days
    - `-3w`: Last 3 weeks
    - `1m`: Next month
    - `5y`: For 5 year plans!

- For more options, run
    ```bash
    todo <command> --help
    ```
