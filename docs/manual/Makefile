clean:
	@rm -f *.aux *.log *.out *.toc *.lol *.lot *.lof *.bbl *.blg *.idx *.ilg *.ind \
		*.nlo *.nls *.nlg *.spl *.synctex.gz *.fdb_latexmk *.fls *.listing *.pdf
pdf:
	@pdflatex -shell-escape -synctext=1 -interaction=nonstopmode \
		-file-line-error manual.tex