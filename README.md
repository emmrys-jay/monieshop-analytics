# Monieshop-Analytics
Monieshop Analytics is an analytics tool for a digital accounting system that tracks
sales transactions. It reads through text transaction files and reports the following metrics:
    -  Highest sales volume in a day
    - Highest sales value in a day
    - Most sold product ID by volume
    - Highest sales staff ID for each month
    - Highest hour of the day by average transaction volume

## Requires
`go 1.18+`

## Setup

- Clone project using:
```bash
    git clone git@github.com:emmrys-jay/monieshop-analytics.git
```

- Enter projects directory using:
```base
    cd monieshop-analytics
```

- Run the program using
```bash
    go build . && ./monieshop --dir "Directory of transaction files"
```
Make sure you change "Directory of transaction files" to your own directory

- You will see results similar to the following:
```bash
    Highest Sales Volume In a Day: 
    Day: 2025-01-28, Sales Volume: 43724

    Highest Sales Value In a Day: 
    Day: 2025-01-28, Sales Value: 2.2228731731999997e+07

    Most Sold Product ID By Volume: 
    Product ID: 583631

    Highest Sales Staff ID for each month: 
    Month: January   StaffID: 5 Volume: 2452
    Month: February  StaffID: 5 Volume: 3884
    Month: March     StaffID: 6 Volume: 2258
    Month: April     StaffID: 2 Volume: 1592
    Month: May       StaffID: 4 Volume: 4229
    Month: June      StaffID: 8 Volume: 1555
    Month: July      StaffID: 8 Volume: 4143
    Month: August    StaffID: 7 Volume: 3379
    Month: September StaffID: 1 Volume: 5451
    Month: October   StaffID: 5 Volume: 2867
    Month: November  StaffID: 5 Volume: 2012
    Month: December  StaffID: 6 Volume: 3886

    Highest Hour of The Day By Average Transaction Volume: 
    Hour: 12
```