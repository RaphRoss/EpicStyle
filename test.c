// =====================================
// Fichier 1: good_example.c (Code propre)
// =====================================
#include <stdio.h>
#include <stdlib.h>

#define MAX_SIZE 100
#define MIN_VALUE 0

/*
** Additionne deux nombres entiers
*/
int add_numbers(int a, int b)
{
	int result;

	result = a + b;
	return (result);
}

/*
** Affiche un message de bienvenue
*/
void print_welcome(void)
{
	printf("Bienvenue dans EpicStyle!\n");
}

/*
** Calcule la moyenne de trois nombres
*/
float calculate_average(float x, float y, float z)
{
	float sum;
	float average;

	sum = x + y + z;
	average = sum / 3.0;
	return (average);
}

int main(void)
{
	int num1;
	int num2;
	int sum;

	num1 = 5;
	num2 = 10;
	sum = add_numbers(num1, num2);
	print_welcome();
	printf("La somme est: %d\n", sum);
	return (0);
}

// =====================================
// Fichier 2: bad_example.c (Avec erreurs)
// =====================================

#include <stdio.h>
#include <stdlib.h>

// Commentaire interdit avec //
#define maxSize 50
#define min_value 0

int globalVar = 42;  // Variable globale non const

int BadFunctionName(int a, int b, int c, int d, int e)  // Trop de paramètres et mauvais nom
{
    int x, y, z;  // Plusieurs variables sur une ligne + espaces au lieu de TAB
    
    for (int i = 0; i < 10; i++) {  // Déclaration dans la boucle
        printf("Ligne très très très très très très très très très très très très longue qui dépasse 80 caractères\n");
    }
    
    int late_declaration = 5;  // Déclaration après du code
    
    return a + b + c + d + e + late_declaration;
}

void function_without_comment(void)
{
	printf("Cette fonction n'a pas de commentaire\n");
}

void very_long_function_that_exceeds_twenty_five_lines(void)
{
	printf("Ligne 1\n");
	printf("Ligne 2\n");
	printf("Ligne 3\n");
	printf("Ligne 4\n");
	printf("Ligne 5\n");
	printf("Ligne 6\n");
	printf("Ligne 7\n");
	printf("Ligne 8\n");
	printf("Ligne 9\n");
	printf("Ligne 10\n");
	printf("Ligne 11\n");
	printf("Ligne 12\n");
	printf("Ligne 13\n");
	printf("Ligne 14\n");
	printf("Ligne 15\n");
	printf("Ligne 16\n");
	printf("Ligne 17\n");
	printf("Ligne 18\n");
	printf("Ligne 19\n");
	printf("Ligne 20\n");
	printf("Ligne 21\n");
	printf("Ligne 22\n");
	printf("Ligne 23\n");
	printf("Ligne 24\n");
	printf("Ligne 25\n");
	printf("Ligne 26\n");
	printf("Ligne 27\n");
}

void fourth_function(void) {}
void fifth_function(void) {}  // Trop de fonctions dans le fichier

int main(void)
{
    return 0;
}


// =====================================
// Fichier 3: BadFileName.c (Nom de fichier incorrect)
// =====================================
#include <stdio.h>

/*
** Fonction principale avec nom de fichier incorrect
*/
int main(void)
{
	printf("Ce fichier a un mauvais nom!\n");
	return (0);
}

// =====================================
// Fichier 4: mixed_errors.c (Erreurs mixtes)
// =====================================
#include <stdio.h>

#define wrong_macro_name 42

void mixedFunction_Name(void)
{
	int var1, var2, var3;
	// Commentaire interdit
	printf("Test\n");
}

int main(void)
{
	return (0);
}

// =====================================
// Fichier 5: header_example.h (Exemple de header)
// =====================================
#ifndef HEADER_EXAMPLE_H
#define HEADER_EXAMPLE_H

#define BUFFER_SIZE 256
#define ERROR_CODE -1

/*
** Structure représentant un point
*/
typedef struct point_s
{
	int x;
	int y;
} point_t;

/*
** Prototypes de fonctions
*/
int calculate_distance(point_t p1, point_t p2);
void print_point(point_t point);
point_t create_point(int x, int y);

#endif /* HEADER_EXAMPLE_H */