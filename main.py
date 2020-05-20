
from graph_tool.all import *
from graph_tool import generation as gn
from random import randrange


class Person:
    def __init__(self, edges=[], infective=False, survived=False, dead=False, infective_epochs=14):
        self.edges = edges
        self.infective = infective
        self.survived = survived
        self.dead = dead
        self.infective_epochs = infective_epochs

def deg_sampler_graph():
    return 1


if __name__ == "__main__":
    '''
    network = []
    
    for i in range(4905854):
        network.append(Person())
    '''
    network_size = 4905854
    connections = 150
    
    #g = Graph(directed=False)
    
    print("Network Graph")
    print("Size: {},\t Edges: {}".format(network_size, connections))
    print("starting...")
    g = gn.random_graph(network_size, deg_sampler_graph, directed=False)
    #vertex = list(g.add_vertex(network_size))
    '''
    g.add_vertex(network_size)
    for v in g.vertices():
        print(v)
    
    for i in range(network_size):
        v = g.add_vertex()
        vertex.append(v)        
        if i%100000 == 0:
            print(i)
    '''
    print("Network Graph Created.")
    '''
    for i in range(network_size):
        for j in range(connections):
            ID = randrange(network_size)
            g.add_edge(vertex[i], vertex[ID])
        if i%100000 == 0:
            print(i)
    '''
    print("Network Graph Edges Added.")
