# deduper
service to identify which contacts are possible matches

### File Path Configuration

This program requires modifying the file path to correctly locate the input file. Make sure to update the path in the code before running the program to avoid file-not-found errors.
To point to your workspace its necessary to change

>	dir  = "/Users/some_user/workspace" 

its same for output of scoring with name match_results.csv
### Similarity Percentage Calculation

The similarity percentage is calculated using the scoreField(name, sourceName, weight) function, where each field contributes a specific percentage to the total 100% score using levenshtein distance for each field. The breakdown is as follows:

| Field       | Weight (%) |
|------------|------------|
| name       | 5%         |
| name1      | 5%         |
| email      | 20%        |
| postalZip  | 30%        |
| address    | 40%        |
| **Total**  | **100%**   |

Each field's similarity score is calculated separately and then combined according to these weightings to determine the overall similarity percentage.

### RUN TESTS 

run on root folder project
>go test -v ./test
