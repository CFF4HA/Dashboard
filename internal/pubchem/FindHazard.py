# imports
import pubchem


if __name__ == "__main__":
    while True:
        chemical = input("Enter a chemical or compound name: ")
        pubchem.Ingredient(chemical)
