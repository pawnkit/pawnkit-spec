new plain[] = "abc";
new marked[] = !"abc";

#assert sizeof plain == 4
#assert sizeof marked == 1

main()
{
    return sizeof plain + sizeof marked;
}
