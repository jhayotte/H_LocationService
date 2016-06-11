FROM scratch
EXPOSE 8001
ADD testlocation /
CMD ["/testlocation"]
