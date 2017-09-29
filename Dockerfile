FROM scratch

COPY _build /

ENTRYPOINT ["/openvidu_tutorial"]
