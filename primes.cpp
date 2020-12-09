
#include<bits/c++.h>
using namespace std;

int main()
{

int var=500;

for(int i=1;i<=500;i++)
{
bool flag=0;
for(int j=2;j<=sqrt(i);j++)
{
if(i%j==0)
{
flag=1;
break;
}
}
if(!flag)
cout<<i<<" ";

}

return 0;
}
