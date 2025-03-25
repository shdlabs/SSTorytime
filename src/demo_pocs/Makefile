#

OBJ=search_broad search_by_node search_clusters search_contrawave_quantum search_simple_text search_by_name search_by_node_text search_contrawave_ambiguous search_interactive search_by_name_unicode search_clusters_functions search_contrawave search_noninteractive search_systematic_knowledge_1 search_systematic_knowledge_2

all: $(OBJ)

% : %.go ../../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go


clean:
	rm -f $(OBJ)
	rm -f *~

