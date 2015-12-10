FROM centos
ADD cuto /cuto/
RUN mkdir /cuto/joblog
RUN mkdir /cuto/log
ENV CUTOROOT /cuto
EXPOSE 2015

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["/cuto/bin/servant"]
