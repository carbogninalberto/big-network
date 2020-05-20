

class Person:
    def __init__(self, edges=[], infective=False, survived=False, dead=False, infective_epochs=14):
        self.edges = edges
        self.infective = infective
        self.survived = survived
        self.dead = dead
        self.infective_epochs = infective_epochs

if __name__ == "__main__":
    network = []
    
    for i in range(4905854):
        network.append(Person())
