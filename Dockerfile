FROM scratch
EXPOSE 8001
ADD H_LocationService /
CMD ["/H_LocationService"]
