FILENAME=cheatsheet

default: $(FILENAME).pdf

$(FILENAME).pdf : $(FILENAME).tex
	pdflatex $(FILENAME)

edit :
	gvim $(FILENAME).tex

spell :
	ispell -t $(FILENAME).tex

view : $(FILENAME).pdf
	evince $(FILENAME).pdf &

clean :
	-rm $(FILENAME).pdf

