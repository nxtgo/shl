#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_NAME_LEN 32

typedef struct
{
    int id;
    char name[MAX_NAME_LEN];
    float grade;
} Student;

void print_student(const Student *s)
{
    if (s == NULL)
    {
        printf("invalid student pointer\n");
        return;
    }
    printf("id: %d\nname: %s\ngrade: %.2f\n\n", s->id, s->name, s->grade);
}

float average_grade(const Student *students, size_t count)
{
    if (count == 0)
        return 0.0f;

    float sum = 0.0f;
    for (size_t i = 0; i < count; i++)
    {
        sum += students[i].grade;
    }
    return sum / (float)count;
}

int main(void)
{
    Student students[3] = {
        {1, "alice", 91.2f}, {2, "bob", 85.5f}, {3, "charlie", 77.8f}};

    printf("=== student records ===\n\n");
    for (int i = 0; i < 3; i++)
    {
        print_student(&students[i]);
    }

    float avg = average_grade(students, 3);
    printf("average grade: %.2f\n", avg);

    // comment
    Student *new_student = malloc(sizeof(Student));
    if (!new_student)
    {
        perror("malloc failed");
        return EXIT_FAILURE;
    }

    new_student->id = 4;
    strcpy(new_student->name, "diana");
    new_student->grade = 88.9f;

    printf("\n=== new student ===\n");
    print_student(new_student);

    free(new_student);
    return EXIT_SUCCESS;
}
