FILENAME=cheatsheet

default: $(FILENAME).pdf

$(FILENAME).pdf : $(FILENAME).tex
	pdflatex $(FILENAME)

edit :
	vim $(FILENAME).tex

spell :
	ispell -t $(FILENAME).tex

view : $(FILENAME).pdf
	evince $(FILENAME).pdf &

clean :
	-rm $(FILENAME).pdf

