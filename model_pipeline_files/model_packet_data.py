import torch
import torch.nn as nn
import torch.nn.functional as F
import torchvision.ops as O



class Model(nn.Module):
  def __init__(self,dimension,num_class):
    super(Model,self).__init__()
    self.dimension=dimension
    self.num_class=num_class
    self.lstm=nn.LSTM(input_size=dimension,hidden_size=dimension,batch_first=True)
    self.conv_layer1=nn.Conv2d(in_channels=1,out_channels=32,kernel_size=3,stride=1,padding="same")
    self.conv_layer2=nn.Conv2d(in_channels=32,out_channels=64,kernel_size=3,stride=1,padding="same")
    self.conv_layer3=nn.Conv2d(in_channels=64,out_channels=64,kernel_size=3,stride=1,padding="same")
    self.bn1 = nn.BatchNorm2d(64)
    self.max_pool1=nn.MaxPool2d(2,stride=2)
    self.se_module=O.SqueezeExcitation(input_channels=64,squeeze_channels=4)
    self.conv_layer4=nn.Conv2d(in_channels=64,out_channels=32,kernel_size=3,stride=1,padding="same")
    self.conv_layer5=nn.Conv2d(in_channels=32,out_channels=32,kernel_size=3,stride=1,padding="same")
    self.conv_layer6=nn.Conv2d(in_channels=32,out_channels=16,kernel_size=3,stride=1,padding="same")
    self.bn2 = nn.BatchNorm2d(16)
    self.max_pool2=nn.MaxPool2d(2,stride=2)
    self.lin1=nn.Linear(dimension*dimension,1024)
    self.relu=nn.ReLU()
    self.dropout=nn.Dropout(0.5)
    self.lin2=nn.Linear(1024,self.num_class)

  def forward(self,x,inference=False):
    x=self.lstm(x)[0]
    x=x.view(-1,1,self.dimension,self.dimension)
    x=self.conv_layer1(x)
    x=self.relu(x)
    x=self.conv_layer2(x)
    x = self.relu(x)
    x=self.conv_layer3(x)
    x = self.bn1(x)
    x = self.relu(x)
    x=self.max_pool1(x)
    x=self.se_module(x)
    x=self.conv_layer4(x)
    x=self.relu(x)
    x=self.conv_layer5(x)
    x = self.relu(x)
    x=self.conv_layer6(x)
    x = self.bn2(x)
    x = self.relu(x)
    x=self.max_pool2(x)
    x=x.view(-1,self.dimension*self.dimension)
    x=self.lin1(x)
    x=self.relu(x)
    if not inference:
      x=self.dropout(x)
    x = self.lin2(x)
    if self.num_class > 1:
      x=torch.softmax(x,dim=-1)
   
    x=torch.sigmoid(x)
    x=torch.round(x)
    return x
# Model(28,1)